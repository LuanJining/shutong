import PageLoading from "./components/PageLoading";
import RouterBeforeEach from "./router/RouterBeforeEach";

const App = () => {
  return (
    <div className="app">
      <PageLoading />
      <RouterBeforeEach />
    </div>
  );
};
export default App;
