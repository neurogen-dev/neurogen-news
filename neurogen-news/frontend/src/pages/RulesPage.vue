<script setup lang="ts">
import { Shield, Check, X, AlertTriangle } from 'lucide-vue-next'
</script>

<template>
  <div class="min-h-screen max-w-3xl mx-auto">
    <div class="bg-white dark:bg-dark-secondary rounded-xl p-8">
      <!-- Header -->
      <div class="flex items-center gap-3 mb-6">
        <Shield class="w-8 h-8 text-primary" />
        <h1 class="text-2xl font-bold text-text-primary dark:text-white">
          Правила сообщества
        </h1>
      </div>
      
      <p class="text-text-secondary mb-8">
        Neurogen.News — это место для конструктивного общения об искусственном интеллекте. 
        Эти правила помогают поддерживать дружелюбную и полезную атмосферу.
      </p>
      
      <!-- Allowed -->
      <section class="mb-8">
        <h2 class="flex items-center gap-2 text-lg font-bold text-success mb-4">
          <Check class="w-5 h-5" />
          Приветствуется
        </h2>
        <ul class="space-y-3">
          <li class="flex gap-3 text-text-secondary">
            <span class="text-success">✓</span>
            Качественный контент о нейросетях и AI
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-success">✓</span>
            Конструктивная критика и обсуждения
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-success">✓</span>
            Помощь новичкам и ответы на вопросы
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-success">✓</span>
            Делиться полезными ресурсами и инструментами
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-success">✓</span>
            Уважительное отношение к другим участникам
          </li>
        </ul>
      </section>
      
      <!-- Not allowed -->
      <section class="mb-8">
        <h2 class="flex items-center gap-2 text-lg font-bold text-error mb-4">
          <X class="w-5 h-5" />
          Запрещено
        </h2>
        <ul class="space-y-3">
          <li class="flex gap-3 text-text-secondary">
            <span class="text-error">✗</span>
            Оскорбления, угрозы и дискриминация
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-error">✗</span>
            Спам и чрезмерная самореклама
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-error">✗</span>
            Плагиат и нарушение авторских прав
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-error">✗</span>
            Распространение вредоносного ПО
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-error">✗</span>
            Публикация незаконного контента
          </li>
          <li class="flex gap-3 text-text-secondary">
            <span class="text-error">✗</span>
            Манипуляции с рейтингом (накрутка)
          </li>
        </ul>
      </section>
      
      <!-- Consequences -->
      <section class="mb-8 p-4 bg-warning/5 border border-warning/20 rounded-xl">
        <h2 class="flex items-center gap-2 text-lg font-bold text-warning mb-3">
          <AlertTriangle class="w-5 h-5" />
          Последствия нарушений
        </h2>
        <ul class="space-y-2 text-text-secondary">
          <li><strong>Первое нарушение:</strong> предупреждение</li>
          <li><strong>Повторное нарушение:</strong> временный бан (1-7 дней)</li>
          <li><strong>Грубые нарушения:</strong> перманентный бан</li>
        </ul>
      </section>
      
      <!-- Content guidelines -->
      <section>
        <h2 class="text-lg font-bold text-text-primary dark:text-white mb-4">
          Рекомендации по контенту
        </h2>
        <div class="space-y-4 text-text-secondary">
          <p>
            <strong>Заголовки:</strong> Должны отражать содержание статьи. 
            Избегайте кликбейта и вводящих в заблуждение формулировок.
          </p>
          <p>
            <strong>Источники:</strong> Указывайте источники информации и 
            давайте ссылки на оригинальные материалы.
          </p>
          <p>
            <strong>NSFW контент:</strong> Должен быть отмечен соответствующим тегом.
          </p>
          <p>
            <strong>Реклама:</strong> Коммерческие публикации должны быть согласованы 
            с администрацией.
          </p>
        </div>
      </section>
      
      <hr class="my-8 border-border dark:border-dark-tertiary" />
      
      <p class="text-sm text-text-tertiary text-center">
        Последнее обновление: декабрь 2024
      </p>
    </div>
  </div>
</template>

