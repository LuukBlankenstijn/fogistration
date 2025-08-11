import { RouterProvider } from "@tanstack/react-router"
import { AuthProvider, useAuth } from "./auth"
import { router } from "./router"
import './styles.css'
import { QueryClientProvider } from "@tanstack/react-query"
import client from "./query/client"
import { ReactQueryDevtools } from "@tanstack/react-query-devtools"

const InnerApp = () => {
  const auth = useAuth()
  return <RouterProvider router={router} context={{ auth }} />
}

const App = () => {
  return (
    <QueryClientProvider client={client}>
      <AuthProvider>
        <InnerApp />
      </AuthProvider>
      <ReactQueryDevtools />
    </QueryClientProvider>
  )
}

export default App;
