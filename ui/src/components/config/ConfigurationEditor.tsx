import React from 'react';
import Editor, { type OnMount } from '@monaco-editor/react';
import { Skeleton } from '@/components/ui/skeleton';

interface ConfigurationEditorProps {
  value: string;
  onChange: (value: string) => void;
  isLoading?: boolean;
}

export const ConfigurationEditor: React.FC<ConfigurationEditorProps> = ({
  value,
  onChange,
  isLoading = false,
}) => {
  const handleEditorDidMount: OnMount = (editor) => {
    editor.focus();
  };

  if (isLoading) {
    return (
      <div className="w-full h-[calc(100vh-280px)] min-h-[500px] rounded-xl border border-border bg-card p-4 space-y-3">
        <Skeleton className="h-6 w-1/3" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-5/6" />
        <Skeleton className="h-4 w-4/6" />
        <Skeleton className="h-[calc(100%-80px)] w-full" />
      </div>
    );
  }

  return (
    <div className="w-full h-[calc(100vh-280px)] min-h-[500px] rounded-xl border border-border overflow-hidden bg-[#1e1e1e] shadow-2xs flex flex-col">
      <div className="bg-[#252526] px-4 py-2.5 border-b border-[#333333] flex items-center justify-between select-none">
        <span className="font-mono text-xs text-muted-foreground uppercase tracking-wider font-semibold">
          sentinel.gateway.yaml
        </span>
        <span className="font-mono text-[11px] text-muted-foreground/80">UTF-8 • YAML</span>
      </div>
      <div className="flex-1 relative">
        <Editor
          height="100%"
          language="yaml"
          theme="vs-dark"
          value={value}
          onChange={(val) => onChange(val ?? '')}
          onMount={handleEditorDidMount}
          options={{
            lineNumbers: 'on',
            automaticLayout: true,
            wordWrap: 'on',
            matchBrackets: 'always',
            quickSuggestions: false,
            minimap: { enabled: false },
            fontSize: 13,
            fontFamily: "'Geist Mono', ui-monospace, SFMono-Regular, monospace",
            scrollBeyondLastLine: false,
            padding: { top: 12 },
          }}
        />
      </div>
    </div>
  );
};
