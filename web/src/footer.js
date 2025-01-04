class CommonFooter extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
<footer class="p-4 bg-white md:p-8 lg:p-10 dark:bg-black mt-4">
    <div class="mx-auto max-w-screen-xl text-center">

        <ul class="flex flex-wrap justify-center items-center mb-6 text-gray-900 dark:text-white">

            <li>
                <a href="/" class="hover:underline me-4 md:me-6">Home</a>
            </li>
            <li>
                <a href="#" class="hover:underline me-4 md:me-6">About</a>
            </li>
            <li>
                <a href="http://github.com/zmaillard/tmbgodt/releases/tag/v1.0.4" class="hover:underline me-4 md:me-6">v1.0.4</a>
            </li>
        </ul>
    </div>
</footer>
        `;
    }
}
window.customElements.define('common-footer', CommonFooter);