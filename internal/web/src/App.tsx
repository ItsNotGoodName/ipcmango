import { createQuery } from "@tanstack/solid-query"
import { createQueryKeyStore } from "@lukemorales/query-key-factory";
import { provideTheme } from "./ui/theme"
import { ErrorBoundary, Suspense } from "solid-js"


export const queries = createQueryKeyStore({
  github: {
    tanstack: null
  }
});

function App() {
  provideTheme()

  const repositoryQuery = createQuery(() => ({
    ...queries.github.tanstack,
    queryFn: async () => {
      const result = await fetch('https://api.github.com/repos/TanStack/query')
      if (!result.ok) throw new Error('Failed to fetch data')
      return result.json()
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
    throwOnError: true, // Throw an error if the query fails
  }))

  return (
    <div>
      <div>Static Content</div>
      {/* An error while fetching will be caught by the ErrorBoundary */}
      <ErrorBoundary fallback={<div>Something went wrong!</div>}>
        {/* Suspense will trigger a loading state while the data is being fetched */}
        <Suspense fallback={<div>Loading...</div>}>
          <div>{repositoryQuery.data?.updated_at}</div>
        </Suspense>
      </ErrorBoundary>
    </div>
  )
}

export default App
