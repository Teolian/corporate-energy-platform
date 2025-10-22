<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <RouterLink to="/" class="text-primary-600 hover:underline mb-4 inline-block">
          ← Back to Home
        </RouterLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Industry Comparison</h1>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="card p-12 text-center">
        <div class="text-6xl mb-4">🏭</div>
        <h2 class="text-2xl font-semibold text-gray-900 dark:text-white mb-4">
          Coming Soon
        </h2>
        <p class="text-gray-600 dark:text-gray-300 mb-6">
          The industry comparison feature is under development
        </p>
        <RouterLink to="/companies" class="btn btn-primary">
          Explore Companies Database
        </RouterLink>
      </div>
    </main>
  </div>
</template>
