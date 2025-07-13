import { Provider } from "@/components/ui/provider";
import { queryClient } from "@/lib/utils";
import { AppRoutes } from "@/routes.tsx";
import { QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <Provider>
        <BrowserRouter>
          <AppRoutes />
        </BrowserRouter>
      </Provider>
    </QueryClientProvider>
  </StrictMode>
);
