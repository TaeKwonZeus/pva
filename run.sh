#!/bin/sh

clear
cd frontend
npm run build
cd ../
go run .
