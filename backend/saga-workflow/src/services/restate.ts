import { tradeWorkflow } from "./saga.ts";
import { serve } from "@restatedev/restate-sdk"
import * as clients from "@restatedev/restate-sdk-clients";

export const restateClient = clients.connect({ url: "http://localhost:8080" }); // the url of your restate server (either cloud or local URL)

serve({
  services: [tradeWorkflow],
  port: 9080
});