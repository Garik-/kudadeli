<script setup lang="ts">
import type { ApexOptions, ApexAxisChartSeries } from 'apexcharts'
import type { ComputedRef } from 'vue'
import { computed } from 'vue'
import { XMarkIcon } from '@heroicons/vue/16/solid'

import VueApexCharts from 'vue3-apexcharts'

import { useExpensesStore } from '@/stories/expensesStore'
import { useFiltersStore } from '@/stories/filtersStore'
import { formatPercent } from '@/utils/formatter'
import { useCategoriesStore } from '@/stories/categoriesStore'

const expensesStore = useExpensesStore()
const filtersStore = useFiltersStore()
const categoriesStore = useCategoriesStore()

function setCategoryFilter(name: string) {
  filtersStore.setFilter('category', name)
}

const props = defineProps({
  onClose: { type: Function, required: true },
})

function handleClose() {
  props.onClose()
}

function getLabelColor(color: string, n: number) {
  const parts = color.split('-')
  parts[parts.length - 1] = n.toString()
  return parts.join('-')
}

function makeOklchTransparent(colorString: string, alpha = 0.7) {
  const oklchRegex = /oklch\((\d+(?:\.\d+)?)% (\d+(?:\.\d+)?) (\d+(?:\.\d+)?)\)/
  const match = colorString.match(oklchRegex)

  if (!match) {
    console.warn('Invalid OKLCH format, returning original color')
    return colorString
  }

  const lightness = parseFloat(match[1])
  const chroma = parseFloat(match[2])
  const hue = parseFloat(match[3])

  return `oklch(${lightness.toFixed(1)}% ${chroma} ${hue} / ${alpha})`
}

function addGlowEffect() {
  const paths = document.querySelectorAll('.apexcharts-pie-series path')
  let cssRules = ''
  paths.forEach((path) => {
    const fillColor = path.getAttribute('fill')!
    const j = path.getAttribute('j')!

    cssRules += `
      .apexcharts-pie-series path[j="${j}"] {
        filter: drop-shadow(0 6px 12px ${makeOklchTransparent(fillColor, 0.35)});
      }

     .apexcharts-pie-series path[j="${j}"]:hover {
        filter: drop-shadow(0 6px 12px ${makeOklchTransparent(fillColor, 0.5)}) brightness(0.96);
      }
    `
  })

  let styleElement = document.getElementById('dynamic-chart-styles')
  if (!styleElement) {
    styleElement = document.createElement('style')
    styleElement.id = 'dynamic-chart-styles'
    document.head.appendChild(styleElement)
  }

  styleElement.textContent = cssRules
}

const chartOptions: ComputedRef<ApexOptions> = computed(() => {
  return {
    chart: {
      id: 'vuechart-example',
      events: {
        mounted: addGlowEffect,
      },
      sparkline: {
        enabled: false,
      },
    },
    grid: {
      padding: {
        top: 30,
        right: 30,
        bottom: 30,
        left: 30,
      },
    },

    legend: {
      show: false,
    },
    stroke: {
      width: 0,
    },
    tooltip: {
      enabled: false,
    },
    colors: expensesStore.groupedByCategory.map(({ color }) => categoriesStore.getHexColor(color)),
    plotOptions: {
      pie: {
        donut: {
          size: '80%',
          background: 'transparent',
        },
        dataLabels: {
          offset: 45,
          minAngleToShowLabel: 20,
        },
        customScale: 1,
        expandOnClick: false,
      },
    },
    dataLabels: {
      enabled: true,
      formatter: function (val: number) {
        return formatPercent(val, 1)
      },
      style: {
        colors: ['var(--color-gray-500)'],
        fontSize: 'var(--text-sm)',
        fontFamily:
          'var(--default-font-family, ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji")',
        fontWeight: '400',
      },
      dropShadow: { enabled: false },
    },
    labels: expensesStore.groupedByCategory.map(({ title }) => title),
  }
})

const series: ComputedRef<ApexAxisChartSeries | ApexNonAxisChartSeries> = computed(() =>
  expensesStore.groupedByCategory.map(({ percent }) => percent),
)
</script>
<template>
  <div class="flex justify-between mb-6">
    <div>
      <div className="text-3xl font-bold">{{ expensesStore.totalAmount }}</div>
      <div className="text-sm">Траты</div>
    </div>
    <div className="w-8 h-8 bg-gray-50 rounded-full p-1 cursor-pointer" @click="handleClose">
      <XMarkIcon class="w-full h-full text-gray-500" />
    </div>
  </div>
  <div>
    <VueApexCharts
      width="100%"
      type="donut"
      :options="chartOptions"
      :series="series"
    ></VueApexCharts>
  </div>
  <div class="flex flex-wrap gap-2 mb-6">
    <div
      v-for="category in expensesStore.groupedByCategory"
      :key="category.amount"
      :class="[
        'flex items-center gap-1 rounded-full p-1 cursor-pointer',
        getLabelColor(category.color, 100),
      ]"
      @click="setCategoryFilter(category.name)"
    >
      <div class="w-8 h-8 rounded-full p-2" :class="category.color">
        <component :is="category.icon" class="w-full h-full text-white" />
      </div>
      <div class="text-sm font-semibold text-gray-700 pr-2">
        {{ category.title }} {{ category.amountFormatted }}
      </div>
    </div>
  </div>
</template>
