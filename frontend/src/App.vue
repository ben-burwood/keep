<template>
    <div class="min-h-screen min-w-screen bg-base-200">
        <ThemeSwitcher class="absolute top-4 right-4 hidden md:block" />

        <div class="flex flex-col items-center p-5 w-full max-w-lg mx-auto">
            <h1 class="text-4xl font-bold hidden md:block">Keep</h1>

            <NewKeep @add="addKeep" class="mt-5 w-full" />
            <div class="divider"></div>

            <div class="flex flex-col gap-2 w-full">
                <Keep
                    v-for="keep in keeps"
                    :key="keep.uuid"
                    :keep
                    @edit="(content) => updateKeep(keep.uuid, content)"
                    @delete="(uuid) => deleteKeep(uuid)"
                />

                <p v-if="keeps.length === 0" class="text-center text-base-content/40 mt-4">No keeps yet</p>
            </div>

            <!-- Error Toast -->
            <div v-if="errorMessage" role="alert" class="alert alert-error fixed bottom-10">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
                {{ errorMessage }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, watch } from "vue";
import ThemeSwitcher from "@/components/ThemeSwitcher.vue";
import { SERVER_URL } from "@/main";
import NewKeep from "@/components/NewKeep.vue";
import Keep from "@/components/Keep.vue";

const errorMessage = ref("");
watch(errorMessage, (newError) => {
    if (newError) {
        setTimeout(() => {
            errorMessage.value = "";
        }, 5000);
    }
});

const keeps = ref<{ uuid: string; content: string }[]>([]);

// API calls
async function fetchKeeps() {
    try {
        const res = await fetch(`${SERVER_URL}/keeps`);
        keeps.value = await res.json();
    } catch (error: any) {
        errorMessage.value = `Error: Fetching Keeps : ${error.message}`;
    }
}
onMounted(fetchKeeps);

async function addKeep(content: string) {
    try {
        await fetch(`${SERVER_URL}/keeps/create`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ content }),
        });
        await fetchKeeps();
    } catch (error: any) {
        errorMessage.value = `Error: Adding Keep : ${error.message}`;
    }
}

async function updateKeep(uuid: string, updatedContent: string) {
    try {
        await fetch(`${SERVER_URL}/keeps/${uuid}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ content: updatedContent }),
        });
        await fetchKeeps();
    } catch (error: any) {
        errorMessage.value = `Error: Updating Keep : ${error.message}`;
    }
}

async function deleteKeep(uuid: string) {
    if (!confirm("Are you sure you want to delete this keep?")) {
        return;
    }
    try {
        await fetch(`${SERVER_URL}/keeps/${uuid}`, { method: "DELETE" });
        await fetchKeeps();
    } catch (error: any) {
        errorMessage.value = `Error: Deleting Keep : ${error.message}`;
    }
}
</script>
