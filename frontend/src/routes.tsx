import LandingPage from "@/pages/LandingPage";
import TaskDetailPage from "@/pages/TaskDetailPage";
import { Route, Routes } from "react-router";

export const AppRoutes = () => (
  <Routes>
    <Route path="/" index element={<LandingPage />} />
    <Route path="/tasks/:id" element={<TaskDetailPage />} />
  </Routes>
);
