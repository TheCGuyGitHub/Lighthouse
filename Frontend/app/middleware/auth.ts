import pb from "~/lib/pocketbase";
export default defineNuxtRouteMiddleware(async (to, from) => {

    if (import.meta.server) {
        return;
    }

    if (!pb.authStore.isValid) {
        console.info('User is not authenticated, redirecting to login page. Authenticated:', pb.authStore.isValid);
        // User is not authenticated, redirect to the login page
        return navigateTo('/login');
    }
});