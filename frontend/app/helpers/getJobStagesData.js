export default function getJobStagesData(stages) {
    const map = {}
    stages.forEach((stage) => {
        map[stage.id] = {
            ...stage,
            startedAt: stage.startedAt ? new Date(stage.startedAt) : stage.startedAt,
            finishedAt: stage.finishedAt ? new Date(stage.finishedAt) : stage.finishedAt,
        };
    });
    return map;
}