import { vitePlugin as remix } from "@remix-run/dev";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [
    remix({
      ignoredRouteFiles: ["**/*.css"],
    

      
    }),
    tsconfigPaths(),
    
  ],
  server: {
    allowedHosts: ["5173-litmuszhang-personalpla-oq7gv3kzeko.ws-eu118.gitpod.io"]
  }
  
});
