export type ValidationStatus = 'ready' | 'valid' | 'error';

export interface ValidationResult {
  status: ValidationStatus;
  line?: number;
  reason?: string;
}
