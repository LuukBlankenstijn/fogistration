import { RouterProvider } from "@tanstack/react-router"
import { router } from "./router"
import './styles.css'
import { QueryClientProvider } from "@tanstack/react-query"
import client from "./query/query-client"
import { ReactQueryDevtools } from "@tanstack/react-query-devtools"

const InnerApp = () => {
  return <RouterProvider router={router} />
}

const App = () => {
  return (
    <QueryClientProvider client={client}>
      <InnerApp />
      <ReactQueryDevtools />
    </QueryClientProvider>
  )
}

export default App;
