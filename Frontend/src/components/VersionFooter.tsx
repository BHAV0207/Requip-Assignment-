import { useEffect, useState } from 'react';
import api from '../services/api';

interface VersionInfo {
  version: string;
  commit: string;
}

// Shows which backend build served this page, so a reviewer can match the UI
// to a commit, a CI run and the image tag Argo CD deployed.
export default function VersionFooter() {
  const [info, setInfo] = useState<VersionInfo | null>(null);

  useEffect(() => {
    api
      .get<VersionInfo>('/version')
      .then((res) => setInfo(res.data))
      .catch(() => setInfo(null));
  }, []);

  return (
    <footer className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-8 text-xs text-slate-400">
      UI {import.meta.env.VITE_APP_VERSION || 'dev'} · API{' '}
      {info ? `${info.version} (${info.commit.slice(0, 7)})` : 'unavailable'}
    </footer>
  );
}
