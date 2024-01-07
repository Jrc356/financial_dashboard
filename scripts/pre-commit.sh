#!/bin/sh

cd server && make test
cd ../client && CI=true npm test