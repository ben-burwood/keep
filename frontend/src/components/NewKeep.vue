<template>
    <div v-if="!newKeepOpen" @click="newKeepOpen = true" class="card bg-base-100 shadow-sm cursor-pointer hover:shadow-md transition-shadow">
        <div class="card-body py-3 px-4">
            <span class="text-base-content/50">New keep...</span>
        </div>
    </div>

    <div v-else class="card bg-base-100 shadow-md">
        <div class="card-body p-4 gap-3">
            <textarea
                ref="newKeepTextarea"
                v-model="newKeepContent"
                placeholder="Write keep (markdown supported)..."
                class="textarea textarea-ghost w-full min-h-32 text-base leading-relaxed resize-y p-0 focus:outline-none focus:bg-transparent"
                @keydown.escape="closeNewKeep"
                @keydown.ctrl.enter="submitNewKeep"
            ></textarea>
            <div class="card-actions justify-end">
                <button class="btn btn-ghost btn-sm" @click="closeNewKeep">Cancel</button>
                <button class="btn btn-accent btn-sm" :disabled="!newKeepContent.trim()" @click="submitNewKeep">Save</button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from "vue";

const newKeepOpen = ref(false);
const newKeepContent = ref("");
const newKeepTextarea = ref<HTMLTextAreaElement | null>(null);

watch(newKeepOpen, async (open) => {
    if (open) {
        await nextTick();
        newKeepTextarea.value?.focus();
    }
});

function closeNewKeep() {
    newKeepOpen.value = false;
    newKeepContent.value = "";
}

const emit = defineEmits(["add"]);

async function submitNewKeep() {
    const content = newKeepContent.value.trim();
    if (!content) return;

    emit("add", content);
    closeNewKeep();
}
</script>
