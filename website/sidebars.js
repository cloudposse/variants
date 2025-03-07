/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

module.exports = {
    tutorials: [
        {
          type: 'category',
          label: 'Tutorials Index Page',
          link: {
            type: 'doc',
            id: 'tutorials/tutorials'
          },
          items: [
            {
                type: 'autogenerated',
                dirName: 'tutorials',
            }
          ]
        }
    ],
    community: [
        {
            type: 'category',
            label: 'Community',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            link: {type: 'doc', id: 'community/community'},
            items: [ 'community/slack', 'community/office-hours' ]
        },
        {
          type: 'category',
          label: 'Contribute',
          className: 'sidebar-title',
          collapsible: false,
          collapsed: false,
          items: [
            {
                type: 'autogenerated',
                dirName: 'contribute',
            }
          ]
        }
    ],
    docs: [

        {
            type: 'category',
            label: 'Get Started',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            items: [
                {
                    type: 'autogenerated',
                    dirName: 'introduction',
                },
                {
                    type: 'category',
                    label: 'Quick Start',
                    collapsible: true,
                    collapsed: false,
                    items: [
                        {
                            type: 'autogenerated',
                            dirName: 'quick-start',
                        },
                    ]
                },
            ]
        },

        ,/*
        {
            type: 'html',
            value: 'Core Concepts',
            className: 'sidebar-title',
        },*/
        {
            type: 'category',
            label: 'Learn Atmos',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            items: [
                {
                    type: 'autogenerated',
                    dirName: 'core-concepts',
                },

                {
                    type: 'category',
                    label: 'Troubleshoot',
                    collapsible: true,
                    collapsed: true,
                    items: [
                        {
                            type: 'autogenerated',
                            dirName: 'troubleshoot',
                        },
                    ]
                }
            ]
        },
        {
            type: 'category',
            label: 'Best Practices',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            link: {type: 'doc', id: 'best-practices/best-practices'},
            items: [
                {
                    type: 'autogenerated',
                    dirName: 'best-practices',
                },
                {
                    type: 'category',
                    label: 'Use Design Patterns',
                    collapsible: true,
                    collapsed: true,
                    link: {type: 'doc', id: 'design-patterns/design-patterns'},
                    items: [
                        {
                            type: 'autogenerated',
                            dirName: 'design-patterns',
                        },
                    ]
                }
            ]
        },
        {
            type: 'category',
            label: 'Cheat Sheets',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            items: [
                {
                    type: 'autogenerated',
                    dirName: 'cheatsheets',
                },
            ]
        }
    ],
    cli: [
        {
            type: 'category',
            label: 'Configuration',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            items: [
                {
                    type: 'autogenerated',
                    dirName: 'cli',
                },
                {
                    type: 'category',
                    label: 'Schemas',
                    collapsible: true,
                    collapsed: true,
                    link: {type: 'doc', id: 'schemas/schemas'},
                    items: [
                        {
                            type: 'autogenerated',
                            dirName: 'schemas',
                        },
                    ]
                }
            ]
        },

        {
            type: 'category',
            label: 'Commands',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            link: {type: 'doc', id: 'cli/commands/commands'},
            items: [
                {
                    type: 'autogenerated',
                    dirName: 'cli/commands',
                },
            ]
        },
        {
            type: 'category',
            label: 'CI/CD (GitOps)',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            items: [
                    {
                        type: 'autogenerated',
                        dirName: 'integrations',
                    }
            ]
        },
        {
            type: 'category',
            label: 'Resources',
            className: 'sidebar-title',
            collapsible: false,
            collapsed: false,
            items: [
                {
                    type: 'autogenerated',
                    dirName: 'reference',
                },
                {
                    type: 'category',
                    label: 'Glossary',
                    collapsible: true,
                    collapsed: true,
                    link: {type: 'doc', id: 'glossary/glossary'},
                    items: [
                        {
                            type: 'autogenerated',
                            dirName: 'glossary',
                        },
                    ]
                }
            ]
        },
    ]
};
