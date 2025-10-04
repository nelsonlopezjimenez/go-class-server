import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  isRouteErrorResponse,
} from 'react-router';

import { Provider as StoreProvider } from 'react-redux';

import { Provider } from '@/components/ui/provider';
import { Layout as AppLayout } from '@/lib/layout';

// fonts
import '@fontsource-variable/plus-jakarta-sans';
import { store } from './api/store';
import {
  ProgressBar,
  ProgressLabel,
  ProgressRoot,
} from './components/ui/progress';

export function Layout({ children }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="UTF-8" />
        <link rel="icon" type="image/svg+xml" href="/assets/favicon.svg" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>CIS Class Server</title>
        <Meta />
        <Links />
      </head>
      <body>
        <StoreProvider store={store}>
          <Provider>
            <AppLayout>{children}</AppLayout>
          </Provider>
          <ScrollRestoration />
          <Scripts />
        </StoreProvider>
      </body>
    </html>
  );
}

export default function Root() {
  return <Outlet />;
}

export function ErrorBoundary({ error }: any) {
  let message = 'Oops!';
  let details = 'An unexpected error occurred.';
  let stack;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? '404' : 'Error';
    details =
      error.status === 404
        ? 'The requested page could not be found.'
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}

export function HydrateFallback() {
  return (
    <ProgressRoot
      size="lg"
      width="full"
      value={null}
      variant="outline"
      colorPalette="blue"
    >
      <ProgressLabel>One moment while your page is loading.</ProgressLabel>
      <ProgressBar />
    </ProgressRoot>
  );
}
