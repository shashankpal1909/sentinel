import React from 'react';
import { CheckCircle2, AlertOctagon, Terminal } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import type { ValidationResult } from '@/types/config';

interface ValidationPanelProps {
  validation: ValidationResult;
}

export const ValidationPanel: React.FC<ValidationPanelProps> = ({ validation }) => {
  return (
    <Card className="w-full">
      <CardHeader className="border-b border-border pb-3">
        <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
          <Terminal className="w-3.5 h-3.5" />
          Syntax Validation
        </CardTitle>
      </CardHeader>
      <CardContent className="pt-4 font-mono text-xs">
        {validation.status === 'ready' && (
          <div className="flex items-center gap-2 text-muted-foreground py-1">
            <span className="w-2 h-2 rounded-full bg-muted-foreground/40 animate-pulse" />
            <span>Ready</span>
          </div>
        )}

        {validation.status === 'valid' && (
          <div className="flex items-center gap-2.5 text-success font-medium py-1">
            <CheckCircle2 className="w-4 h-4 shrink-0" />
            <span>YAML syntax valid</span>
          </div>
        )}

        {validation.status === 'error' && (
          <div className="space-y-2 py-1">
            <div className="flex items-center gap-2.5 text-error font-semibold">
              <AlertOctagon className="w-4 h-4 shrink-0" />
              <span>YAML syntax error</span>
            </div>
            {validation.line && (
              <div className="text-[11px] text-muted-foreground">
                Line <span className="font-bold text-foreground">{validation.line}</span>
              </div>
            )}
            {validation.reason && (
              <div className="p-2.5 rounded-lg bg-error/10 border border-error/20 text-error/90 font-mono text-[11px] leading-relaxed break-all">
                {validation.reason}
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
