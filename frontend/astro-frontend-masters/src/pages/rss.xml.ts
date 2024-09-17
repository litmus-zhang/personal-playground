import rss from '@astrojs/rss'
import { getCollection } from 'astro:content'
import type { AstroConfig } from 'astro'
import sanitizeHtml from 'sanitize-html'
import MarkdownIt from 'markdown-it'

const parser = new MarkdownIt();

export async function get(context: AstroConfig){
    const blogs = await getCollection('blog');
    return rss({
        title: 'The Sandwch Blog',
        description: 'The latest news from the Sandwich blog',
        site: context.site ?? 'https://sandwich.dev',
        items: blogs.map((blog) => {
            return {
                title: blog.data.title,
                pubDate: new Date(blog.data.date),
                link: `/blog/${blog.slug}`,
                author: "Sandwich",
                content: sanitizeHtml(parser.render(blog.data.description)),


            }
        })
    })
}