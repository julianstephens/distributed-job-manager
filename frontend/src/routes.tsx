import { Route, Routes } from "react-router";
import App from "@/App";

export const AppRoutes = () => (
  <Routes>
    <Route path="/" index element={<App />} />
  </Routes>
);
