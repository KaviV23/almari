import { Button } from '@/components/ui/button'
import { Breadcrumb, BreadcrumbList, BreadcrumbItem, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import { Home, List, LayoutGrid, ChevronRight } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function Breadcrumbs() {
  const { breadcrumbs, loadDir, navigateToBreadcrumb, viewMode, setViewMode } = useFileManagerContext()

  return (
    <div className="bg-gray-800 border border-gray-700 px-4 py-3 rounded-xl flex items-center justify-between shadow">
      <div className="flex items-center space-x-2 overflow-x-auto text-sm">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <Button
                variant="ghost"
                size="sm"
                className="text-indigo-400 hover:text-indigo-300 font-semibold h-auto p-0"
                onClick={() => loadDir('/')}
              >
                <Home className="w-4 h-4 mr-1" />
                <span>/</span>
              </Button>
            </BreadcrumbItem>
            {breadcrumbs.map((segment, index) => (
              <BreadcrumbItem key={index}>
                <BreadcrumbSeparator>
                  <ChevronRight className="w-3 h-3 text-gray-500" />
                </BreadcrumbSeparator>
                <Button
                  variant="ghost"
                  size="sm"
                  className="hover:text-indigo-300 text-gray-300 font-medium h-auto p-0"
                  onClick={() => navigateToBreadcrumb(index)}
                >
                  {segment}
                </Button>
              </BreadcrumbItem>
            ))}
          </BreadcrumbList>
        </Breadcrumb>
      </div>
      <div className="flex items-center space-x-4">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="icon" className={viewMode === 'list' ? 'bg-indigo-600 text-white' : ''}>
              <List className="w-4 h-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => setViewMode('list')}>
              <List className="mr-2 h-4 w-4" />
              List View
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => setViewMode('icon')}>
              <LayoutGrid className="mr-2 h-4 w-4" />
              Icon View
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <div className="text-xs text-gray-400">
          <span id="file-count">0 items</span>
        </div>
      </div>
    </div>
  )
}
