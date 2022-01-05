/**
 * Copyright 2020 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package token ...
package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// IsTokenValid checks whether the provided token is valid or not
func IsTokenValid(logger *zap.Logger, tokenString string) bool {
	logger.Error("Validating token")

	token, err := parseToken(tokenString)
	if err != nil {
		logger.Error("Error parsing token", zap.Error(err))
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := claims.Valid(); err != nil {
			logger.Error("Token is invalid", zap.Error(err))
			return false
		}
		logger.Info("Token is valid")
		return true
	}

	logger.Error("Unable to fetch token claims")
	return false
}

// FetchTokenLifeTime fetches token life time of the token
func FetchTokenLifeTime(logger *zap.Logger, tokenString string, tokenExpirydiff uint64) (uint64, error) {
	logger.Info("Fetching token life time")
	var tokenLifeTime uint64

	token, err := parseToken(tokenString)
	if err != nil {
		logger.Error("Error parsing token", zap.Error(err))
		return tokenLifeTime, errors.New("error parsing the token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := claims.Valid(); err != nil {
			logger.Error("Token is invalid", zap.Error(err))
			return tokenLifeTime, err
		}
		currentTime := time.Now().Unix()
		var expiryTime interface{}
		if expiryTime, ok = claims["exp"]; !ok {
			logger.Error("Unable to find expiry time of token")
			return tokenLifeTime, errors.New("unable to find expiry time of token")
		}
		tokenLifeTime = uint64(expiryTime.(float64)) - uint64(currentTime)
		if tokenLifeTime < tokenExpirydiff {
			logger.Error("Token life time is less than expected", zap.Uint64("Expected token expiry diff", tokenExpirydiff), zap.Uint64("Token life time", tokenLifeTime))
			return tokenLifeTime, errors.New("token life time is less than expected")
		}
		logger.Info("Successfully fetched token life time")
		return tokenLifeTime, nil
	}
	logger.Error("Unable to fetch token claims")
	return tokenLifeTime, errors.New("unable to fetch token claims")
}

// parseToken parses token string to jwt token
func parseToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, errors.New("empty token string")
	}
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	return token, err
}
