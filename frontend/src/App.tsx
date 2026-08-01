import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { FileManager } from '@/components/file-manager'
import { TooltipProvider } from '@/components/ui/tooltip'
import './App.css'

const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <TooltipProvider>
        <FileManager />
      </TooltipProvider>
    ),
  },
])

function App() {
  return <RouterProvider router={router} />
}

export default App