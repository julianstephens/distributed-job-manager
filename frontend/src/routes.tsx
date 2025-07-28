import JobDashboardPage from "@/pages/JobDashboardPage";
import JobDetailPage from "@/pages/JobDetailPage";
import LandingPage from "@/pages/LandingPage";
import { withAuthenticationRequired } from "@auth0/auth0-react";
import type React from "react";
import { Route, Routes } from "react-router";

const ProtectedRoute = ({
  component,
  ...args
}: {
  component: React.ComponentType;
}) => {
  const Component = withAuthenticationRequired(component, args);
  return <Component />;
};

export const AppRoutes = () => (
  <Routes>
    <Route path="/" index element={<LandingPage />} />
    <Route
      path="/jobs"
      index
      element={<ProtectedRoute component={JobDashboardPage} />}
    />
    <Route
      path="/jobs/:id"
      index
      element={<ProtectedRoute component={JobDetailPage} />}
    />
  </Routes>
);
