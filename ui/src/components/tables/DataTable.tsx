import React from 'react';
import { Database } from 'lucide-react';
import { cn } from '@/lib/utils';
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from '@/components/ui/table';

export interface Column<T> {
  header: string;
  accessorKey?: keyof T;
  cell?: (item: T) => React.ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  keyExtractor: (item: T, index: number) => string;
  emptyMessage?: string;
  className?: string;
}

export function DataTable<T>({
  columns,
  data,
  keyExtractor,
  emptyMessage = 'No records found.',
  className,
}: DataTableProps<T>): React.ReactElement {
  return (
    <div
      className={cn(
        'w-full rounded-xl border border-border bg-card overflow-hidden shadow-2xs transition-colors duration-150',
        className
      )}
    >
      <Table>
        <TableHeader className="bg-secondary/60">
          <TableRow className="hover:bg-transparent border-b border-border">
            {columns.map((col, idx) => (
              <TableHead
                key={idx}
                className={cn(
                  'h-11 px-5 font-mono text-[11px] font-semibold uppercase tracking-wider text-muted-foreground select-none',
                  col.className
                )}
              >
                {col.header}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody className="divide-y divide-border font-mono text-xs">
          {data.length === 0 ? (
            <TableRow className="hover:bg-transparent">
              <TableCell
                colSpan={columns.length}
                className="px-5 py-14 text-center text-muted-foreground font-sans"
              >
                <div className="flex flex-col items-center justify-center gap-2.5">
                  <Database className="w-6 h-6 text-muted-foreground/60" />
                  <span className="text-xs font-medium">{emptyMessage}</span>
                </div>
              </TableCell>
            </TableRow>
          ) : (
            data.map((item, rowIdx) => (
              <TableRow
                key={keyExtractor(item, rowIdx)}
                className="hover:bg-secondary/50 transition-colors duration-150 group border-b border-border last:border-0"
              >
                {columns.map((col, colIdx) => {
                  let content: React.ReactNode;
                  if (col.cell) {
                    content = col.cell(item);
                  } else if (col.accessorKey) {
                    content = String(item[col.accessorKey] ?? '');
                  }
                  return (
                    <TableCell
                      key={colIdx}
                      className={cn(
                        'px-5 py-3.5 text-foreground/90 group-hover:text-foreground leading-normal',
                        col.className
                      )}
                    >
                      {content}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}
