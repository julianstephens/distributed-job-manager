import App from "@/App";
import TaskDetailPage from "@/pages/TaskDetailPage";
import { Route, Routes } from "react-router";

export const AppRoutes = () => (
  <Routes>
    <Route path="/" index element={<App />} />
    <Route path="/tasks/:id" element={<TaskDetailPage />} />
  </Routes>
);
