/**
 * Parses a label_values(metric, label) query string.
 * Returns the metric and label, or null if the format doesn't match.
 */
export function parseLabelValuesQuery(query: string): { metric: string; label: string } | null {
  const match = query.match(/^label_values\(\s*([^,]+?)\s*,\s*([^)]+?)\s*\)$/);
  if (!match) return null;
  return { metric: match[1], label: match[2] };
}

/**
 * Substitutes template variables in a string.
 * Replaces ${var} and $var patterns with their values.
 * Processes longer variable names first to avoid prefix collisions.
 */
export function substituteVariables(template: string, variables: Record<string, string>): string {
  if (!template || Object.keys(variables).length === 0) return template;

  let result = template;

  // Sort variable names by length descending to avoid prefix collisions
  const names = Object.keys(variables).sort((a, b) => b.length - a.length);

  for (const name of names) {
    const value = variables[name];
    // Replace ${var} form
    result = result.split(`\${${name}}`).join(value);
    // Replace $var form (only when not followed by word characters to avoid partial matches)
    result = result.replace(new RegExp(`\\$${name}(?![a-zA-Z0-9_])`, 'g'), value);
  }

  return result;
}
