package main

import (
	utils "LoveLetterProject/internal"
	"LoveLetterProject/internal/database"
	"LoveLetterProject/pkg/authorization"
	"LoveLetterProject/pkg/authorization/endpoints"
	"LoveLetterProject/pkg/authorization/middleware"
	"LoveLetterProject/pkg/authorization/transport"
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/juju/ratelimit"
	"github.com/oklog/oklog/pkg/group"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// schema for user table
const userSchema = `
		create table if not exists users (
			id 		   Varchar(36) not null,
			email 	   Varchar(100) not null unique,
			username   Varchar(225),
			password   Varchar(225) not null,
			tokenhash  Varchar(15) not null,
			verified   Boolean default false,
		    banned     Boolean default false,
			createdat  Timestamp not null,
			updatedat  Timestamp not null,
			Primary Key (id)
		);
`

// schema for verification table
const verificationSchema = `
		create table if not exists verifications (
			email 		Varchar(100) not null,
			code  		Varchar(10) not null,
			expiresat 	Timestamp not null,
			type        Varchar(10) not null,		
			Primary Key (email),
			Constraint fk_user_email Foreign Key(email) References users(email)
				On Delete Cascade On Update Cascade
		)
`

// schema for profile table
const profileSchema = `
		create table if not exists profiles (
			id 		   Varchar(36) not null,
			userid 	Varchar(36) not null,
		    email 	   Varchar(100) not null,
			firstname   Varchar(225) default '',
			lastname    Varchar(225) default '',
			avatarurl  Varchar(255) default '',
			phone      Varchar(25) default '',
			street     Varchar(255) default '',
			city       Varchar(255) default '',
			state      Varchar(10) default '',
			zipcode   Varchar(5) default '',
			country    Varchar(255) default '',
			createdat  Timestamp not null,
			updatedat  Timestamp not null,
			Primary Key (id),
			Constraint fk_user_id Foreign Key(userid) References users(id)
				On Delete Cascade On Update Cascade
		)
`

// schema for securityuser table
const securityUserSchema = `
		create table if not exists passworusers (
			id 		   Varchar(36) not null,
			userid 	Varchar(36) not null,
			password   Varchar(225) not null,	
			createdat  Timestamp not null,
			updatedat  Timestamp not null,
			Primary Key (id),
			Constraint fk_user_id Foreign Key(userid) References users(id)
				On Delete Cascade On Update Cascade
		)
`

// schema for limit table
const limitSchema = `
		create table if not exists limits (
			id 		   Varchar(36) not null,
			userid 	Varchar(36) not null,
			numofsendmail 	   Int default 0,
			numofchangepassword Int default 0,
		    numoflogin 	   Int default 0,
			createdat  Timestamp not null,
			updatedat  Timestamp not null,
			Primary Key (id),
			Constraint fk_user_id Foreign Key(userid) References users(id)
				On Delete Cascade On Update Cascade
		)
`

// schema for multiratio table
const multiratioSchema = `
		create table if not exists multiratios (
		waterratio int default 0,
		lightratio int default 0,
		seedratio int default 0,
		unique  (waterratio, lightratio, seedratio)
		)
`

// schema for earnscore table
const earnscoreSchema = `
		create table if not exists earnscores (			
			userid 	Varchar(36) not null,
			waterscore 	   Int default 0,
			lightscore Int default 0,
			seedscore 	   Int default 0,
			createdat  Timestamp not null,
			updatedat  Timestamp not null,
		    			Primary Key (userid),
			Constraint fk_user_id Foreign Key(userid) References users(id)
				On Delete Cascade On Update Cascade
		)
`

func main() {
	logger := utils.NewLogger()

	// quynhlx change config with multi environments
	configs := utils.NewConfigurations(logger, utils.DeployLocal)
	//configs := utils.NewConfigurations(logger, utils.DeployStage)
	//configs := utils.NewConfigurations(logger, utils.DeployProd)
	// validator contains all the methods that are need to validate the user json in request
	validator := database.NewValidation()
	// create a new connection to the postgres db store
	db, err := database.NewConnection(configs, logger)
	defer db.Close()
	if err != nil {
		logger.Error("unable to connect to db", "error", err)
		return
	}
	// creation of user table.
	db.MustExec(userSchema)
	db.MustExec(verificationSchema)
	db.MustExec(profileSchema)
	db.MustExec(securityUserSchema)
	db.MustExec(limitSchema)
	db.MustExec(multiratioSchema)
	db.MustExec(earnscoreSchema)

	// repository contains all the methods that interact with DB to perform CURD operations for user.
	repository := database.NewPostgresRepository(db, logger)
	// mailService contains the utility methods to send an email
	mailService := authorization.NewSGMailService(logger, configs)
	// authService contains all methods that help in authorizing a user request
	auth := middleware.NewAuthService(logger, configs)

	// Reset limit data for users.
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("00:01").Do(func() {
		logger.Info("Clear limit data for users after 1 day at 00:01.")
		ctx := context.Background()
		err := repository.ClearAllLimitData(ctx)
		if err != nil {
			logger.Error("Error clearing limit data", "error", err)
		}
	})
	s.StartAsync()

	// Create rate limiter for users.
	rlBucket := ratelimit.NewBucket(1*time.Second, 5)

	var (
		httpAddr    = net.JoinHostPort("localhost", configs.HttpPort)
		service     = authorization.NewUserService(logger, configs, repository, mailService, auth)
		eps         = endpoints.NewEndpointSet(service, auth, repository, logger, validator, rlBucket)
		httpHandler = transport.NewHTTPHandler(eps)
	)

	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			//logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			logger.Error("could not start the server", "error", err)
			os.Exit(1)
		}
		g.Add(func() error {
			//logger.Log("transport", "HTTP", "addr", httpAddr)
			logger.Info("starting the server at port", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			logger.Error("could not start the server", "error", err)
			httpListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Info("exit", g.Run())
}
