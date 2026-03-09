<template>
    <div class="border border-base-300/50 shadow-sm rounded-lg bg-base-100 hover:bg-primary/10">
        <div class="card-body p-4">
            <div v-if="!isEditing" @click="startEdit" class="cursor-pointer">
                <div class="prose prose-sm break-words" v-html="renderedContent"></div>
            </div>

            <div v-else class="flex flex-col gap-3">
                <textarea
                    ref="editTextarea"
                    v-model="editContent"
                    class="textarea textarea-ghost w-full min-h-24 text-base leading-relaxed resize-y p-0 focus:outline-none focus:bg-transparent"
                    @keydown.ctrl.enter="submitEdit"
                    @keydown.escape="cancelEdit"
                ></textarea>
                <div class="card-actions justify-between">
                    <button class="btn btn-ghost btn-sm btn-error" @click="emit('delete', keep.uuid)">Delete</button>
                    <div class="flex gap-2">
                        <button class="btn btn-ghost btn-sm" @click="cancelEdit">Cancel</button>
                        <button class="btn btn-primary btn-sm" :disabled="!editContent.trim()" @click="submitEdit">Save</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from "vue";
import MarkdownIt from "markdown-it";

const md = new MarkdownIt({ linkify: true });

interface Keep {
    uuid: string;
    content: string;
}

const props = defineProps<{ keep: Keep }>();
const renderedContent = computed(() => md.render(props.keep.content));

const isEditing = ref(false);
const editContent = ref("");
const editTextarea = ref<HTMLTextAreaElement | null>(null);

const emit = defineEmits(["edit", "delete"]);

function startEdit() {
    isEditing.value = true;
    editContent.value = props.keep.content;
    nextTick(() => editTextarea.value?.focus());
}

function cancelEdit() {
    isEditing.value = false;
    editContent.value = "";
}

async function submitEdit() {
    const content = editContent.value.trim();
    if (!content) return;

    emit("edit", content);
    cancelEdit();
}
</script>
