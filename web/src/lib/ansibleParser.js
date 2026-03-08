/**
 * Parse Ansible playbook output into structured format
 */
export class AnsibleOutputParser {
  constructor() {
    this.plays = [];
    this.currentPlay = null;
    this.currentTask = null;
    this.facts = {};
    this.handlers = [];
    this.variables = {};
  }

  parse(logs) {
    this.reset();

    logs.forEach((log) => {
      const output = log.output || '';

      // Parse PLAY
      const playMatch = output.match(/^PLAY\s+\[([^\]]+)\]/);
      if (playMatch) {
        this.startPlay(playMatch[1], log.time);
        return;
      }

      // Parse TASK
      const taskMatch = output.match(/^TASK\s+\[([^\]]+)\]/);
      if (taskMatch) {
        this.startTask(taskMatch[1], log.time);
        return;
      }

      // Parse task result
      const resultMatch = output.match(/^(ok|changed|failed|skipped|unreachable):\s*\[([^\]]+)\]\s*(.*)/);
      if (resultMatch) {
        this.addTaskResult(resultMatch[1], resultMatch[2], resultMatch[3], log.time);
        return;
      }

      // Parse facts gathering
      if (output.includes('GATHERING FACTS')) {
        this.currentTask = {
          name: 'Gathering Facts',
          type: 'setup',
          status: 'running',
          startTime: log.time,
        };
        return;
      }

      // Parse handlers
      const handlerMatch = output.match(/^RUNNING HANDLER\s+\[([^\]]+)\]/);
      if (handlerMatch) {
        this.startHandler(handlerMatch[1], log.time);
        return;
      }

      // Parse variables (vars_prompt, set_fact, etc.)
      const varMatch = output.match(/^ok:\s*\[([^\]]+)\]\s*=>\s*\{[^}]*"ansible_facts"[^}]*\}/);
      if (varMatch) {
        this.parseFacts(output);
        return;
      }

      // Add output to current task
      if (this.currentTask) {
        if (!this.currentTask.output) {
          this.currentTask.output = [];
        }
        this.currentTask.output.push(log);
      }
    });

    // Finalize last play and task
    if (this.currentTask && this.currentPlay) {
      this.currentPlay.tasks.push(this.currentTask);
    }
    if (this.currentPlay) {
      this.plays.push(this.currentPlay);
    }

    return {
      plays: this.plays,
      facts: this.facts,
      handlers: this.handlers,
      variables: this.variables,
    };
  }

  startPlay(name, time) {
    if (this.currentPlay) {
      if (this.currentTask) {
        this.currentPlay.tasks.push(this.currentTask);
        this.currentTask = null;
      }
      this.plays.push(this.currentPlay);
    }

    this.currentPlay = {
      name,
      startTime: time,
      tasks: [],
      hosts: new Set(),
      status: 'running',
    };
    this.currentTask = null;
  }

  startTask(name, time) {
    if (this.currentTask && this.currentPlay) {
      this.currentPlay.tasks.push(this.currentTask);
    }

    this.currentTask = {
      name,
      startTime: time,
      status: 'running',
      results: [],
      output: [],
    };
  }

  addTaskResult(status, host, message, time) {
    if (!this.currentTask) return;

    if (this.currentPlay) {
      this.currentPlay.hosts.add(host);
    }

    this.currentTask.results.push({
      host,
      status,
      message,
      time,
    });

    // Update task status
    if (status === 'failed' || status === 'unreachable') {
      this.currentTask.status = 'failed';
    } else if (status === 'changed' && this.currentTask.status !== 'failed') {
      this.currentTask.status = 'changed';
    } else if (this.currentTask.status === 'running') {
      this.currentTask.status = 'ok';
    }
  }

  startHandler(name, time) {
    this.handlers.push({
      name,
      startTime: time,
      status: 'running',
      results: [],
    });
  }

  parseFacts(output) {
    try {
      // Extract JSON from output
      const jsonMatch = output.match(/\{.*"ansible_facts".*\}/);
      if (jsonMatch) {
        const facts = JSON.parse(jsonMatch[0]);
        Object.assign(this.facts, facts.ansible_facts || {});
      }
    } catch (e) {
      console.warn('Failed to parse facts:', e);
    }
  }

  reset() {
    this.plays = [];
    this.currentPlay = null;
    this.currentTask = null;
    this.facts = {};
    this.handlers = [];
    this.variables = {};
  }
}

export const ansibleParser = new AnsibleOutputParser();
