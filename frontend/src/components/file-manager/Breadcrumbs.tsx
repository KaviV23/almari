import { Button } from '@/components/ui/button'
import { Breadcrumb, BreadcrumbList, BreadcrumbItem, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import { Home, List, LayoutGrid } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function Breadcrumbs() {
  const { breadcrumbs, loadDir, navigateToBreadcrumb, viewMode, setViewMode } = useFileManagerContext()

  return (
    <div className="bg-card border border-border px-4 py-0 rounded-xl flex items-center justify-between shadow">
      <div className="flex items-center space-x-2 overflow-x-auto text-sm">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <Button
                variant="ghost"
                size="sm"
                className="px-1 cursor-pointer"
                onClick={() => loadDir("/")}
              >
                <Home size={"1rem"} className="" />
              </Button>
            </BreadcrumbItem>
            {breadcrumbs.map((segment, index) => (
              <BreadcrumbItem key={index}>
                <BreadcrumbSeparator></BreadcrumbSeparator>
                <Button
                  variant="ghost"
                  size="sm"
                  className="hover:text-card-foreground text-gray-300 font-normal h-auto p-0"
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
            <Button
              variant="ghost"
              size="sm"
              className={"font-normal text-card-foreground p-1"}
            >
              {viewMode === "list" ? (
                <>
                  <List className="h-4 w-4" />
                  List View
                </>
              ) : (
                <>
                  <LayoutGrid className="h-4 w-4" />
                  Icon View
                </>
              )}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => setViewMode("list")} className='text-muted-foreground'>
              <List className="mr-2 h-4 w-4" />
              List View
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => setViewMode("icon")} className='text-muted-foreground'>
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
  );
}
