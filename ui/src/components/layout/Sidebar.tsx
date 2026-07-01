import React from 'react';
import { NavLink } from 'react-router-dom';
import { LayoutDashboard, Cpu, Server, Route, Layers, ShieldCheck } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';

const navItems = [
  { name: 'Dashboard', path: '/', icon: LayoutDashboard },
  { name: 'Runtime', path: '/runtime', icon: Cpu },
  { name: 'Services', path: '/services', icon: Server },
  { name: 'Routes', path: '/routes', icon: Route },
  { name: 'Backends', path: '/backends', icon: Layers },
];

export const Sidebar: React.FC = () => {
  return (
    <aside
      className="w-64 bg-card border-r border-border flex flex-col shrink-0 transition-colors duration-150 select-none"
      aria-label="Main Navigation"
    >
      {/* Brand Header */}
      <div className="h-14 flex items-center gap-3 px-5 border-b border-border">
        <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-primary text-primary-foreground shadow-2xs">
          <ShieldCheck className="w-4 h-4" />
        </div>
        <div className="flex flex-col">
          <span className="font-mono font-bold tracking-wider text-sm uppercase text-foreground leading-none">
            Sentinel
          </span>
          <span className="text-[11px] text-muted-foreground tracking-tight mt-0.5">
            Gateway Control
          </span>
        </div>
      </div>

      {/* Nav List */}
      <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        <div className="px-3 pb-2 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
          Infrastructure
        </div>
        {navItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.path}
              to={item.path}
              end={item.path === '/'}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 px-3 py-2.5 text-xs font-medium rounded-lg transition-all duration-150 relative group focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
                  isActive
                    ? 'bg-secondary text-foreground font-semibold shadow-2xs'
                    : 'text-muted-foreground hover:bg-secondary/60 hover:text-foreground'
                )
              }
            >
              {({ isActive }) => (
                <>
                  {isActive && (
                    <span className="absolute left-0 top-2 bottom-2 w-1 rounded-r bg-primary" />
                  )}
                  <Icon
                    className={cn(
                      'w-4 h-4 shrink-0 transition-colors duration-150',
                      isActive
                        ? 'text-primary'
                        : 'text-muted-foreground group-hover:text-foreground'
                    )}
                  />
                  <span>{item.name}</span>
                </>
              )}
            </NavLink>
          );
        })}
      </nav>

      {/* Footer Info */}
      <div className="p-4 border-t border-border bg-secondary/30">
        <div className="flex items-center justify-between text-[11px] font-mono text-muted-foreground">
          <span>API Epoch</span>
          <Badge variant="outline" className="text-[10px] font-semibold tracking-wider px-1.5 py-0">
            READ-ONLY
          </Badge>
        </div>
      </div>
    </aside>
  );
};
