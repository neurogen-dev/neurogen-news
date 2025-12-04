<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'
import Underline from '@tiptap/extension-underline'
import TextAlign from '@tiptap/extension-text-align'
import Highlight from '@tiptap/extension-highlight'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { common, createLowlight } from 'lowlight'
import {
  Bold,
  Italic,
  Underline as UnderlineIcon,
  Strikethrough,
  Code,
  Quote,
  List,
  ListOrdered,
  Image as ImageIcon,
  Link as LinkIcon,
  AlignLeft,
  AlignCenter,
  AlignRight,
  Highlighter,
  Heading1,
  Heading2,
  Heading3,
  Undo,
  Redo,
  CodeSquare,
  Minus,
} from 'lucide-vue-next'

interface Props {
  modelValue?: string
  placeholder?: string
  autofocus?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: 'Начните писать...',
  autofocus: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const lowlight = createLowlight(common)

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      codeBlock: false,
    }),
    Placeholder.configure({
      placeholder: props.placeholder,
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        class: 'text-primary hover:underline',
      },
    }),
    Image.configure({
      HTMLAttributes: {
        class: 'rounded-lg max-w-full',
      },
    }),
    Underline,
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    Highlight.configure({
      multicolor: true,
    }),
    CodeBlockLowlight.configure({
      lowlight,
    }),
  ],
  editorProps: {
    attributes: {
      class: 'prose-article focus:outline-none min-h-[300px]',
    },
  },
  onUpdate: () => {
    emit('update:modelValue', editor.value?.getHTML() || '')
  },
  autofocus: props.autofocus,
})

// Watch for external content changes
watch(() => props.modelValue, (newValue) => {
  const isSame = editor.value?.getHTML() === newValue
  if (!isSame) {
    editor.value?.commands.setContent(newValue, false)
  }
})

// Link modal
const showLinkModal = ref(false)
const linkUrl = ref('')

const setLink = () => {
  if (linkUrl.value) {
    editor.value?.chain().focus().extendMarkRange('link').setLink({ href: linkUrl.value }).run()
  } else {
    editor.value?.chain().focus().extendMarkRange('link').unsetLink().run()
  }
  showLinkModal.value = false
  linkUrl.value = ''
}

// Image upload
const addImage = () => {
  const url = window.prompt('URL изображения')
  if (url) {
    editor.value?.chain().focus().setImage({ src: url }).run()
  }
}

onBeforeUnmount(() => {
  editor.value?.destroy()
})

const toolbarGroups = [
  {
    name: 'history',
    items: [
      { icon: Undo, action: () => editor.value?.chain().focus().undo().run(), active: false, title: 'Отменить' },
      { icon: Redo, action: () => editor.value?.chain().focus().redo().run(), active: false, title: 'Повторить' },
    ],
  },
  {
    name: 'headings',
    items: [
      { icon: Heading1, action: () => editor.value?.chain().focus().toggleHeading({ level: 1 }).run(), active: () => editor.value?.isActive('heading', { level: 1 }), title: 'Заголовок 1' },
      { icon: Heading2, action: () => editor.value?.chain().focus().toggleHeading({ level: 2 }).run(), active: () => editor.value?.isActive('heading', { level: 2 }), title: 'Заголовок 2' },
      { icon: Heading3, action: () => editor.value?.chain().focus().toggleHeading({ level: 3 }).run(), active: () => editor.value?.isActive('heading', { level: 3 }), title: 'Заголовок 3' },
    ],
  },
  {
    name: 'formatting',
    items: [
      { icon: Bold, action: () => editor.value?.chain().focus().toggleBold().run(), active: () => editor.value?.isActive('bold'), title: 'Жирный' },
      { icon: Italic, action: () => editor.value?.chain().focus().toggleItalic().run(), active: () => editor.value?.isActive('italic'), title: 'Курсив' },
      { icon: UnderlineIcon, action: () => editor.value?.chain().focus().toggleUnderline().run(), active: () => editor.value?.isActive('underline'), title: 'Подчёркнутый' },
      { icon: Strikethrough, action: () => editor.value?.chain().focus().toggleStrike().run(), active: () => editor.value?.isActive('strike'), title: 'Зачёркнутый' },
      { icon: Highlighter, action: () => editor.value?.chain().focus().toggleHighlight().run(), active: () => editor.value?.isActive('highlight'), title: 'Выделение' },
      { icon: Code, action: () => editor.value?.chain().focus().toggleCode().run(), active: () => editor.value?.isActive('code'), title: 'Код' },
    ],
  },
  {
    name: 'alignment',
    items: [
      { icon: AlignLeft, action: () => editor.value?.chain().focus().setTextAlign('left').run(), active: () => editor.value?.isActive({ textAlign: 'left' }), title: 'По левому краю' },
      { icon: AlignCenter, action: () => editor.value?.chain().focus().setTextAlign('center').run(), active: () => editor.value?.isActive({ textAlign: 'center' }), title: 'По центру' },
      { icon: AlignRight, action: () => editor.value?.chain().focus().setTextAlign('right').run(), active: () => editor.value?.isActive({ textAlign: 'right' }), title: 'По правому краю' },
    ],
  },
  {
    name: 'blocks',
    items: [
      { icon: List, action: () => editor.value?.chain().focus().toggleBulletList().run(), active: () => editor.value?.isActive('bulletList'), title: 'Маркированный список' },
      { icon: ListOrdered, action: () => editor.value?.chain().focus().toggleOrderedList().run(), active: () => editor.value?.isActive('orderedList'), title: 'Нумерованный список' },
      { icon: Quote, action: () => editor.value?.chain().focus().toggleBlockquote().run(), active: () => editor.value?.isActive('blockquote'), title: 'Цитата' },
      { icon: CodeSquare, action: () => editor.value?.chain().focus().toggleCodeBlock().run(), active: () => editor.value?.isActive('codeBlock'), title: 'Блок кода' },
      { icon: Minus, action: () => editor.value?.chain().focus().setHorizontalRule().run(), active: false, title: 'Разделитель' },
    ],
  },
  {
    name: 'media',
    items: [
      { icon: LinkIcon, action: () => { linkUrl.value = editor.value?.getAttributes('link').href || ''; showLinkModal.value = true }, active: () => editor.value?.isActive('link'), title: 'Ссылка' },
      { icon: ImageIcon, action: addImage, active: false, title: 'Изображение' },
    ],
  },
]
</script>

<template>
  <div class="border border-border dark:border-dark-tertiary rounded-xl overflow-hidden bg-white dark:bg-dark-secondary">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center gap-1 p-2 border-b border-border dark:border-dark-tertiary bg-background-secondary dark:bg-dark-tertiary">
      <template v-for="(group, groupIndex) in toolbarGroups" :key="group.name">
        <div 
          v-if="groupIndex > 0" 
          class="w-px h-6 bg-border dark:bg-dark-tertiary mx-1"
        />
        <button
          v-for="item in group.items"
          :key="item.title"
          @click="item.action"
          :title="item.title"
          class="p-1.5 rounded transition-colors"
          :class="[
            typeof item.active === 'function' && item.active()
              ? 'bg-primary/10 text-primary'
              : 'text-text-secondary hover:text-text-primary hover:bg-background-tertiary dark:hover:bg-dark'
          ]"
        >
          <component :is="item.icon" class="w-4 h-4" />
        </button>
      </template>
    </div>
    
    <!-- Editor -->
    <EditorContent :editor="editor" class="px-4 py-3" />
    
    <!-- Link Modal -->
    <Teleport to="body">
      <div 
        v-if="showLinkModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <div 
          class="absolute inset-0 bg-black/50"
          @click="showLinkModal = false"
        />
        <div class="relative bg-white dark:bg-dark-secondary rounded-xl p-6 w-full max-w-md shadow-modal">
          <h3 class="text-lg font-bold mb-4 text-text-primary dark:text-white">
            Вставить ссылку
          </h3>
          <input
            v-model="linkUrl"
            type="url"
            placeholder="https://example.com"
            class="w-full px-4 py-2 border border-border dark:border-dark-tertiary rounded-lg bg-white dark:bg-dark-tertiary text-text-primary dark:text-white mb-4"
            @keydown.enter="setLink"
          />
          <div class="flex justify-end gap-2">
            <button
              @click="showLinkModal = false"
              class="px-4 py-2 text-text-secondary hover:text-text-primary transition-colors"
            >
              Отмена
            </button>
            <button
              @click="setLink"
              class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-600 transition-colors"
            >
              Вставить
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style>
/* TipTap editor styles */
.ProseMirror {
  @apply min-h-[300px] outline-none;
}

.ProseMirror p.is-editor-empty:first-child::before {
  @apply text-text-tertiary pointer-events-none float-left h-0;
  content: attr(data-placeholder);
}

.ProseMirror pre {
  @apply bg-dark dark:bg-dark-tertiary rounded-lg p-4 overflow-x-auto my-4;
}

.ProseMirror pre code {
  @apply bg-transparent p-0 text-sm;
}

.ProseMirror img {
  @apply rounded-lg max-w-full my-4;
}

.ProseMirror hr {
  @apply my-6 border-border dark:border-dark-tertiary;
}

.ProseMirror blockquote {
  @apply pl-4 border-l-4 border-primary/30 italic text-text-secondary my-4;
}

.ProseMirror ul,
.ProseMirror ol {
  @apply pl-6 my-4;
}

.ProseMirror ul {
  @apply list-disc;
}

.ProseMirror ol {
  @apply list-decimal;
}

.ProseMirror li {
  @apply mb-1;
}

.ProseMirror mark {
  @apply bg-yellow-200 dark:bg-yellow-500/30 px-0.5;
}
</style>

