<template>
    <div class="border border-base-300/50 shadow-sm rounded-lg bg-base-100 hover:bg-primary/10">
        <div class="card-body p-4">
            <div v-if="!isEditing" @click="handleContentClick" class="cursor-pointer">
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

// Plugin: render `- [ ] …` and `- [x] …` as interactive checkboxes
md.core.ruler.after("inline", "task_lists", (state) => {
    let taskIndex = 0;
    const tokens = state.tokens;

    for (let i = 1; i < tokens.length; i++) {
        const token = tokens[i];
        if (!token || token.type !== "inline") continue;
        if (!/^\[[ xX]\][ \u00a0]/.test(token.content)) continue;

        // Support both tight lists (inline directly after list_item_open)
        // and loose lists (inline after paragraph_open after list_item_open)
        const prev = tokens[i - 1];
        const prevPrev = i >= 2 ? tokens[i - 2] : undefined;

        const prevIsListItem = prev?.type === "list_item_open";
        const prevIsParagraph =
            prev?.type === "paragraph_open" && prevPrev?.type === "list_item_open";

        if (!prevIsListItem && !prevIsParagraph) continue;

        const checked = /^\[[xX]\]/.test(token.content);
        const currentIndex = taskIndex++;

        // Inject checkbox as an html_inline token at the start of children
        const checkboxToken = new state.Token("html_inline", "", 0);
        checkboxToken.content = `<input type="checkbox" data-task-index="${currentIndex}"${checked ? " checked" : ""}>`;

        const children = token.children;
        if (children && children.length > 0) {
            const first = children[0];
            if (first && first.type === "text") {
                // Strip the leading task marker (e.g. `[ ] ` or `[x] `) using a regex
                first.content = first.content.replace(/^\[[xX ]\][ \u00a0]/, "");
            }
            children.unshift(checkboxToken);
        }

        // Add class to list item for styling (remove bullet)
        const listItemToken = prevIsListItem ? prev : prevPrev;
        listItemToken?.attrSet("class", "task-list-item");
    }
});

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

function handleContentClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (
        target instanceof HTMLInputElement &&
        target.type === "checkbox" &&
        target.hasAttribute("data-task-index")
    ) {
        // Prevent the browser from toggling the checkbox — we handle state via emit
        event.preventDefault();
        const index = parseInt(target.getAttribute("data-task-index")!);
        if (isNaN(index)) return;
        toggleTask(index);
        return;
    }
    startEdit();
}

function toggleTask(taskIndex: number) {
    let count = 0;
    const newContent = props.keep.content.replace(
        /^([ \t]*[-*+] )\[([xX ])\][ \u00a0]/gm,
        (match, prefix, status) => {
            if (count++ === taskIndex) {
                return `${prefix}[${status.toLowerCase() === "x" ? " " : "x"}] `;
            }
            return match;
        },
    );
    emit("edit", newContent);
}

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

<style scoped>
/* Remove the default list bullet from task list items */
:deep(.task-list-item) {
    list-style: none;
    margin-left: -1.5em;
    padding-left: 0;
}

:deep(.task-list-item input[type="checkbox"]) {
    margin-right: 0.4em;
    cursor: pointer;
}
</style>
