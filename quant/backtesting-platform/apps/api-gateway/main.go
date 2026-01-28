import "google/api/annotations.proto";

rpc RunBacktest(RunBacktestRequest)
returns (RunBacktestResponse) {
  option (google.api.http) = {
    post: "/v1/backtests/run"
    body: "*"
  };
}