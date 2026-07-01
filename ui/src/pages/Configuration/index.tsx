import React, { useState, useEffect, useCallback } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import * as yaml from 'js-yaml';
import { toast } from 'sonner';
import { PageHeader } from '@/components/layout/PageHeader';
import { ConfigurationEditor } from '@/components/config/ConfigurationEditor';
import { RuntimeSummary } from '@/components/config/RuntimeSummary';
import { ValidationPanel } from '@/components/config/ValidationPanel';
import { ConfigToolbar } from '@/components/config/ConfigToolbar';
import { getConfig, applyConfig, reloadConfig } from '@/api/config';
import type { ValidationResult } from '@/types/config';

export const ConfigurationPage: React.FC = () => {
  const queryClient = useQueryClient();

  const [initialYaml, setInitialYaml] = useState<string>('');
  const [currentYaml, setCurrentYaml] = useState<string>('');
  const [isLoadingInitial, setIsLoadingInitial] = useState<boolean>(true);
  const [isApplying, setIsApplying] = useState<boolean>(false);
  const [isReloading, setIsReloading] = useState<boolean>(false);
  const [validation, setValidation] = useState<ValidationResult>({ status: 'ready' });

  const isDirty = Boolean(initialYaml && currentYaml && initialYaml.trim() !== currentYaml.trim());

  // Load initial config
  const fetchInitialConfig = useCallback(async () => {
    setIsLoadingInitial(true);
    try {
      const yamlStr = await getConfig();
      setInitialYaml(yamlStr);
      setCurrentYaml(yamlStr);
      setValidation({ status: 'ready' });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to fetch configuration';
      toast.error(`Load error: ${msg}`);
    } finally {
      setIsLoadingInitial(false);
    }
  }, []);

  useEffect(() => {
    fetchInitialConfig();
  }, [fetchInitialConfig]);

  // Prompt before unload if dirty
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (isDirty) {
        e.preventDefault();
        e.returnValue = '';
      }
    };
    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [isDirty]);

  // Client-side YAML validation
  const validateYaml = useCallback((yamlStr: string): ValidationResult => {
    try {
      yaml.load(yamlStr);
      return { status: 'valid' };
    } catch (err: unknown) {
      if (typeof err === 'object' && err !== null && 'mark' in err) {
        const mark = err as {
          mark: { line: number; column: number };
          reason?: string;
          message?: string;
        };
        return {
          status: 'error',
          line: mark.mark.line + 1,
          reason: mark.reason || mark.message || 'YAML syntax error',
        };
      }
      return {
        status: 'error',
        reason: err instanceof Error ? err.message : 'Invalid YAML syntax',
      };
    }
  }, []);

  const handleValidate = useCallback(() => {
    const res = validateYaml(currentYaml);
    setValidation(res);
    if (res.status === 'valid') {
      toast.success('✔ YAML syntax valid');
    } else {
      toast.error(`Invalid YAML syntax on line ${res.line || 'unknown'}`);
    }
  }, [currentYaml, validateYaml]);

  // Apply configuration
  const handleApply = useCallback(async () => {
    if (!isDirty) return;

    // Syntax check first
    const syntaxCheck = validateYaml(currentYaml);
    setValidation(syntaxCheck);
    if (syntaxCheck.status === 'error') {
      toast.error(`Cannot apply: YAML syntax error on line ${syntaxCheck.line || 'unknown'}`);
      return;
    }

    setIsApplying(true);
    try {
      const res = await applyConfig(currentYaml);
      toast.success(res.message || 'Configuration applied successfully.');
      setInitialYaml(currentYaml);
      setValidation({ status: 'valid' });
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['runtime'] }),
        queryClient.invalidateQueries({ queryKey: ['backends'] }),
        queryClient.invalidateQueries({ queryKey: ['services'] }),
        queryClient.invalidateQueries({ queryKey: ['routes'] }),
      ]);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Configuration rejected';
      toast.error(`Configuration rejected: ${msg}`);
      setValidation({
        status: 'error',
        reason: msg,
      });
    } finally {
      setIsApplying(false);
    }
  }, [isDirty, currentYaml, validateYaml, queryClient]);

  // Reload configuration from disk
  const handleReload = useCallback(async () => {
    setIsReloading(true);
    try {
      const res = await reloadConfig();
      toast.success(res.message || 'Configuration reloaded from disk.');
      const freshYaml = await getConfig();
      setInitialYaml(freshYaml);
      setCurrentYaml(freshYaml);
      setValidation({ status: 'ready' });
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['runtime'] }),
        queryClient.invalidateQueries({ queryKey: ['backends'] }),
        queryClient.invalidateQueries({ queryKey: ['services'] }),
        queryClient.invalidateQueries({ queryKey: ['routes'] }),
      ]);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to reload configuration';
      toast.error(`Reload failed: ${msg}`);
    } finally {
      setIsReloading(false);
    }
  }, [queryClient]);

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 's') {
        e.preventDefault();
        if (isDirty && !isApplying && !isReloading) {
          handleApply();
        }
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isDirty, isApplying, isReloading, handleApply]);

  const handleEditorChange = (val: string) => {
    setCurrentYaml(val);
    if (validation.status !== 'ready') {
      setValidation({ status: 'ready' });
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Configuration Center"
        description="Live operational workflow for gateway routes, listeners, and backend clusters"
      />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
        {/* Editor column (2 spans) */}
        <div className="lg:col-span-2 space-y-4">
          <ConfigurationEditor
            value={currentYaml}
            onChange={handleEditorChange}
            isLoading={isLoadingInitial}
          />
          <ConfigToolbar
            isDirty={isDirty}
            isValidating={false}
            isApplying={isApplying}
            isReloading={isReloading}
            onValidate={handleValidate}
            onReload={handleReload}
            onApply={handleApply}
          />
        </div>

        {/* Sidebar panels (1 span) */}
        <div className="space-y-6">
          <RuntimeSummary />
          <ValidationPanel validation={validation} />
        </div>
      </div>
    </div>
  );
};
