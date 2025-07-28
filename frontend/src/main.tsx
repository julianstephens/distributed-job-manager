import { Provider } from "@/components/ui/provider";
import { Toaster } from "@/components/ui/toaster";
import { queryClient } from "@/lib/utils";
import { AppRoutes } from "@/routes";
import {
  Auth0Provider,
  User,
  type AppState,
  type Auth0ContextInterface,
  type Auth0ProviderOptions,
} from "@auth0/auth0-react";
import { QueryClientProvider } from "@tanstack/react-query";
import React, { StrictMode, type PropsWithChildren } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, useNavigate } from "react-router";
import "./main.css";

const Auth0ProviderWithRedirectCallback = ({
  children,
  context,
  ...props
}: PropsWithChildren<Omit<Auth0ProviderOptions, "context">> & {
  context?: React.Context<Auth0ContextInterface<User>>;
}) => {
  const navigate = useNavigate();
  const onRedirectCallback = (appState?: AppState) => {
    navigate(appState?.returnTo || window.location.pathname);
  };

  return (
    <Auth0Provider
      onRedirectCallback={onRedirectCallback}
      context={context}
      {...props}
    >
      {children}
    </Auth0Provider>
  );
};

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Provider>
      <BrowserRouter>
        <Auth0ProviderWithRedirectCallback
          domain={import.meta.env.VITE_AUTH0_DOMAIN}
          clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
          authorizationParams={{ redirect_uri: "http://localhost:5173/jobs" }}
        >
          <QueryClientProvider client={queryClient}>
            <AppRoutes />
            <Toaster />
          </QueryClientProvider>
        </Auth0ProviderWithRedirectCallback>
      </BrowserRouter>
    </Provider>
  </StrictMode>
);
