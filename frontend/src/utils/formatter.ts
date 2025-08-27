export function formatPrice(amount: number): string {
  return amount.toLocaleString('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  })
}

export function formatPercent(n: number, fixed = 2): string {
  const rounded = Number(n.toFixed(fixed))
  return (rounded % 1 === 0 ? rounded.toString() : rounded.toFixed(fixed)) + '%'
}

export function capitalizeFirstLetter(str: string) {
  if (typeof str !== 'string' || str.length === 0) return str

  // Удаляем ведущие пробелы, но запоминаем их для восстановления
  const leadingSpaces = str.match(/^\s*/)?.[0] ?? ''
  const trimmed = str.trimStart()

  if (trimmed.length === 0) return str // строка из одних пробелов

  // Берём первый символ с учётом суррогатных пар (UTF-16)
  const firstChar = [...trimmed][0]
  const firstCharUpper = firstChar.toLocaleUpperCase()

  // Остальная часть строки после первого символа
  const rest = [...trimmed].slice(1).join('')

  return leadingSpaces + firstCharUpper + rest
}
