import { Layout } from "@/components/layout";
import { TaskTable } from "@/components/TaskTable";

const App = () => {
  return (
    <Layout title="Task List">
      <TaskTable />
    </Layout>
  );
};

export default App;
