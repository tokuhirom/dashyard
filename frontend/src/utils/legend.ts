type LegendFunc = (value: string, arg: string) => string;

const legendFuncs: Record<string, LegendFunc> = {
  trunc(value, arg) {
    const n = parseInt(arg, 10);
    if (isNaN(n) || n <= 0 || value.length <= n) return value;
    return value.slice(0, n) + '...';
  },
  suffix(value, arg) {
    const n = parseInt(arg, 10);
    if (isNaN(n) || n <= 0 || value.length <= n) return value;
    return '...' + value.slice(-n);
  },
  upper(value) {
    return value.toUpperCase();
  },
  lower(value) {
    return value.toLowerCase();
  },
  replace(value, arg) {
    // replace("old","new")
    const m = arg.match(/^"([^"]*)"\s*,\s*"([^"]*)"$/);
    if (!m) return value;
    return value.split(m[1]).join(m[2]);
  },
};

// Parse "funcName(arg)" or "funcName" and return [funcName, arg], or null.
function parseFuncCall(expr: string): [string, string] | null {
  const m = expr.match(/^(\w+)\(([^)]*)\)$/);
  if (m) return [m[1], m[2]];
  // No-arg form: "funcName"
  const m2 = expr.match(/^(\w+)$/);
  if (m2 && legendFuncs[m2[1]]) return [m2[1], ''];
  return null;
}

function applyPipe(value: string, pipeExpr: string): string {
  const parts = pipeExpr.split('|');
  let result = value;
  for (let i = 1; i < parts.length; i++) {
    const parsed = parseFuncCall(parts[i].trim());
    if (!parsed) continue;
    const [funcName, arg] = parsed;
    const fn = legendFuncs[funcName];
    if (fn) {
      result = fn(result, arg);
    }
  }
  return result;
}

export function buildLabel(metric: Record<string, string>, legend?: string): string {
  if (legend) {
    return legend.replace(/\{([^}]+)\}/g, (_, expr: string) => {
      const labelName = expr.split('|')[0].trim();
      const value = metric[labelName] ?? '';
      if (!expr.includes('|')) return value;
      return applyPipe(value, expr);
    });
  }
  const entries = Object.entries(metric).filter(([k]) => k !== '__name__');
  if (entries.length === 0) {
    return metric['__name__'] || 'value';
  }
  return entries.map(([k, v]) => `${k}="${v}"`).join(', ');
}
