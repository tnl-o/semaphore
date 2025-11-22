export default async function copyToClipboard(content, message) {
  try {
    const el = document.createElement('textarea');
    el.value = content;
    el.setAttribute('readonly', '');
    el.style.position = 'absolute';
    el.style.left = '-9999px';
    document.body.appendChild(el);
    const selected = document.getSelection().rangeCount > 0
      ? document.getSelection().getRangeAt(0) : false;
    el.select();
    const successful = document.execCommand('copy');
    document.body.removeChild(el);
    if (selected) {
      document.getSelection().removeAllRanges();
      document.getSelection().addRange(selected);
    }

    if (!successful) {
      throw new Error('Fallback copy failed');
    }

    const { toast } = await import('@/lib/toast');
    toast.success(message);
  } catch (e) {
    const { toast } = await import('@/lib/toast');
    toast.error(`Can't copy to clipboard: ${e.message}`);
  }
}
