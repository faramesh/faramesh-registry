export default function DisabledRegistryUI() {
  return (
    <main
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: '#111827',
        color: '#9ca3af',
        fontFamily: 'system-ui, sans-serif',
        padding: '2rem',
      }}
    >
      <div>
        <h1 style={{ color: '#f3f4f6', fontSize: '1.25rem', marginBottom: '0.75rem' }}>
          Registry web UI disabled
        </h1>
        <p style={{ maxWidth: '32rem', lineHeight: 1.6, margin: 0 }}>
          Use the GitHub catalog and the Faramesh CLI. See{' '}
          <a href="https://github.com/faramesh/faramesh-registry" style={{ color: '#93c5fd' }}>
            github.com/faramesh/faramesh-registry
          </a>
          .
        </p>
        <p style={{ marginTop: '1rem', fontSize: '0.875rem', color: '#6b7280' }}>
          <code style={{ color: '#d1d5db' }}>faramesh registry list</code>
        </p>
      </div>
    </main>
  );
}
