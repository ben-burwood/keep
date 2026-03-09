<template>
    <div class="min-h-screen min-w-screen bg-base-200">
        <ThemeSwitcher class="absolute top-4 right-4 hidden md:block" />

        <div class="flex flex-col items-center justify-center p-5 w-full max-w-lg mx-auto">
            <h1 class="text-4xl font-bold hidden md:block">Keep</h1>

            <div class="divider"></div>

            <div></div>

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
import { ref, computed, onMounted, watch } from "vue";
import ThemeSwitcher from "@/components/ThemeSwitcher.vue";
import { SERVER_URL } from "@/main";

const errorMessage = ref("");
// clear error after 5 seconds
watch(errorMessage, (newError) => {
    if (newError) {
        setTimeout(() => {
            errorMessage.value = "";
        }, 5000);
    }
});

const keeps = ref<{ uuid: string; content: string }[]>([]);

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
        const res = await fetch(`${SERVER_URL}/keeps/create`, {
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
        const res = await fetch(`${SERVER_URL}/keeps/${uuid}`, {
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
        const res = await fetch(`${SERVER_URL}/keeps/${uuid}`, { method: "DELETE" });
        await fetchKeeps();
    } catch (error: any) {
        errorMessage.value = `Error: Deleting Keep : ${error.message}`;
    }
}
</script>
