'use client';

import React, { createContext, useContext, ReactNode, useMemo } from 'react';
import ApiClient from '@/lib/api-client';

const ApiContext = createContext<ApiClient | undefined>(undefined);

export function ApiProvider({ children }: { children: ReactNode }) {
  const apiClient = useMemo(() => new ApiClient(), []);

  return (
    <ApiContext.Provider value={apiClient}>
      {children}
    </ApiContext.Provider>
  );
}

export function useApi(): ApiClient {
  const context = useContext(ApiContext);
  if (context === undefined) {
    throw new Error('useApi must be used within an ApiProvider');
  }
  return context;
}
