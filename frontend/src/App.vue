<script setup lang="ts">
import ExpenseList from './components/ExpenseList.vue';
</script>

<template>
  <div class="page">
    <div class="page-wrapper">
      <div class="page-body">
        <div class="container-xl">
          <div class="row row-deck row-cards">
            <div class="col-12">
              <ExpenseList></ExpenseList>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
