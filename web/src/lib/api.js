import apiClient from '@/lib/apiClient';

export async function loadProjectResources(name) {
  return (await apiClient({
    method: 'get',
    url: `/api/project/${this.projectId}/${name}`,
    responseType: 'json',
  })).data;
}

export async function test() {
  return null;
}
