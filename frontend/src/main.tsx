import { Provider } from "@/components/ui/provider";
import { AppRoutes } from "@/routes.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";

const queryClient = new QueryClient();
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
