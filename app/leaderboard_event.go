package main

type LeaderboardEvent struct {
	FunctionARN  string `json:"function_arn"`
	FunctionName string `json:"function_name"`
	AccountID    string `json:"account_id"`
	Points       int    `json:"points"`
}
