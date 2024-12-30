export default function getVisualizerDataObject(data) {
    if (!data) return {};
    
    const idToObjectMap = {};
    // Init
    data.forEach(item => {
        idToObjectMap[item.id] = {
            ...item,
            dependencies: [],
            parents: [],
            bloodline: [],
            row: 0,
            startedAt: item.startedAt ? new Date(item.startedAt) : item.startedAt,
            finishedAt: item.finishedAt ? new Date(item.finishedAt) : item.finishedAt,
        };
        item.dependencyDepth = 0;
    });

    // Fill dependencies
    data.forEach(item => {
        if (item.dependencies) {
            item.dependencies.forEach(depId => {
                const parent = idToObjectMap[item.id];
                const dependency = idToObjectMap[depId];

                parent.dependencies.push(dependency);
                dependency.parents.push(parent);

                if (!parent.bloodline.includes(dependency)) {
                    parent.bloodline.push(dependency);
                }
                if (!dependency.bloodline.includes(parent)) {
                    dependency.bloodline.push(parent);
                }
            });
        }
    });
    
    return idToObjectMap;
}

export function getVisualizerDataByLevel(data) {
    if (!data) return {};

    const values = Object.values(data);


    function assignLevels(nodes) {
        const inDegree = {}; // Dependencies
        const levels = {};
        const queue = [];

        nodes.forEach((node) => {
            inDegree[node.id] = node.dependencies.length;
            if (inDegree[node.id] === 0) {
                queue.push(node); // No dependencies => starting nodes
            }
        });

        let level = 0;
        // BFS
        while (queue.length) {
            const currentLevel = [];
            const nextQueue = [];

            queue.forEach((node) => {
                node.dependencyDepth = level;
                currentLevel.push(node);

                node.parents.forEach((parent) => {
                    inDegree[parent.id]--; // Decrease amount of dependencies
                    if (inDegree[parent.id] === 0) {
                        nextQueue.push(parent);
                    }
                });
            });

            levels[level] = currentLevel;
            queue.splice(0, queue.length, ...nextQueue);
            level++;
        }

        return levels;
    }

    // The algorithm
    function minimizeCrossings(levels) {
        const layerKeys = Object.keys(levels).sort((a, b) => +a - +b);

        const firstLayer = levels[layerKeys[0]];
        firstLayer.forEach((node, idx) => {
            node.row = idx;
        });

        for (let i = 1; i < layerKeys.length; i++) {
            const currentLayer = levels[layerKeys[i]];

            // Calculate barycenter for nodes in the current layer
            const barycenter = currentLayer.map((node) => {
                const connectedNodes = node.dependencies.map((dep) => dep.row);
                const avgPosition = connectedNodes.reduce((sum, pos) => sum + pos, 0) / connectedNodes.length;
                return { node, avgPosition };
            });

            // Sort current layer by barycenter
            barycenter.sort((a, b) => a.avgPosition - b.avgPosition);
            currentLayer.forEach((entry, i) => {
                entry.row = i; // Update row position
            });
        }

        return levels;
    }

    // Step 3: Generate Leveled Data
    const levels = assignLevels(values);
    const optimizedLevels = minimizeCrossings(levels);

    const leveledAsIds = {};
    for (const [layer, nodes] of Object.entries(optimizedLevels)) {
        leveledAsIds[layer] = nodes.map((node) => node.id);
    }

    return leveledAsIds; 
}
